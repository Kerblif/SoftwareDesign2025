using HSEFinance.Lib.Application.Facades;
using Spectre.Console;
using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;

namespace HSEFinance.ConsoleApp
{
    public class AccountManagerFacade
    {
        private readonly ImportExportFacade<BankAccount> _accountImportExportFacade;
        private readonly IAccountRepository _accountRepository;

        public AccountManagerFacade(IAccountRepository accountRepository)
        {
            _accountImportExportFacade = new ImportExportFacade<BankAccount>();
            _accountRepository = new AccountRepositoryProxy(accountRepository);
        }

        public void ShowMenu()
        {
            while (true)
            {
                // Главное меню управления счетами
                var choice = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("[green]Управление счетами[/]")
                        .AddChoices("Добавить счет", "Показать все счета", "Удалить счет", "Редактировать счет", "Пересчет", "Импорт", "Экспорт", "Назад"));

                switch (choice)
                {
                    case "Добавить счет":
                        AddAccount();
                        break;

                    case "Показать все счета":
                        ShowAllAccounts();
                        break;

                    case "Удалить счет":
                        DeleteAccount();
                        break;

                    case "Редактировать счет":
                        EditAccount();
                        break;
                    
                    case "Импорт":
                        Import();
                        break;

                    case "Экспорт":
                        Export();
                        break;

                    case "Пересчет":
                        RecalculateAccount();
                        break;
                    
                    case "Назад":
                        return; // Возврат в предыдущее меню
                }
            }
        }

        private void AddAccount()
        {
            try
            {
                // Ввод данных нового счета
                var name = AnsiConsole.Ask<string>("Введите название счета:");

                // Сохранение счета в репозитории
                _accountRepository.CreateBankAccount(name);

                AnsiConsole.MarkupLine($"[green]Счет '{Markup.Escape(name)}' успешно добавлен![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка добавления счета: {Markup.Escape(ex.Message)}[/]");
            }
        }

        private void ShowAllAccounts()
        {
            try
            {
                // Получение списка всех счетов из репозитория
                var accounts = _accountRepository.GetAllBankAccounts().ToList();

                if (accounts.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Список счетов пуст.[/]");
                    return;
                }

                // Отображение данных о счетах
                AnsiConsole.MarkupLine("[green]Список счетов:[/]");

                var table = new Table()
                    .AddColumn("ID")
                    .AddColumn("Название")
                    .AddColumn("Баланс");

                foreach (var account in accounts)
                {
                    table.AddRow(account.Id.ToString(), account.Name, account.Balance.ToString("C"));
                }

                AnsiConsole.Render(table);
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка отображения счетов: {Markup.Escape(ex.Message)}[/]");
            }
        }

        private void DeleteAccount()
        {
            try
            {
                var accounts = _accountRepository.GetAllBankAccounts().ToList();

                if (accounts.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Нет счетов для удаления.[/]");
                    return;
                }

                var accountToDelete = AnsiConsole.Prompt(
                    new SelectionPrompt<BankAccount>()
                        .Title("Выберите счет для удаления:")
                        .AddChoices(accounts));

                _accountRepository.DeleteBankAccount(accountToDelete.Id);

                AnsiConsole.MarkupLine($"[green]Счет '{Markup.Escape(accountToDelete.Name)}' успешно удален![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка удаления счета: {ex.Message}[/]");
            }
        }
        
        private void EditAccount()
        {
            try
            {
                var accounts = _accountRepository.GetAllBankAccounts().ToList();
        
                if (accounts.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Нет счетов для редактирования.[/]");
                    return;
                }
        
                var accountToEdit = AnsiConsole.Prompt(
                    new SelectionPrompt<BankAccount>()
                        .Title("Выберите счет для редактирования:")
                        .AddChoices(accounts));
        
                // Prompt for a new name
                var newName = AnsiConsole.Ask<string>($"Введите новое название для счета (текущее: {accountToEdit.Name}):");
                accountToEdit.Name = newName;
        
                // Prompt for a new balance
                var newBalance = AnsiConsole.Prompt(
                    new TextPrompt<decimal>($"Введите новый баланс для счета (текущий: {accountToEdit.Balance:C}):"));
        
                accountToEdit.Balance = newBalance;
        
                _accountRepository.UpdateBankAccount(accountToEdit);
        
                AnsiConsole.MarkupLine($"[green]Счет '{Markup.Escape(accountToEdit.Name)}' успешно обновлен![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка редактирования счета: {ex.Message}[/]");
            }
        }
        
        private string PromptFormatSelection()
        {
            return AnsiConsole.Prompt(
                new SelectionPrompt<string>()
                    .Title("Выберите формат файла ([green]json[/], [green]csv[/]):")
                    .AddChoices(new[] { "json", "csv" }));
        }

        private void Import()
        {
            try
            {
                var filePath = AnsiConsole.Ask<string>("Введите путь к файлу для импорта:");
        
                var accounts = _accountImportExportFacade.Import(PromptFormatSelection(), filePath);

                if (accounts == null)
                {
                    AnsiConsole.MarkupLine("[red]Файл не содержит никаких допустимых счетов.[/]");
                    return;
                }
        
                foreach (var account in accounts)
                {
                    try
                    {
                        _accountRepository.UploadBankAccount(account);
                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine($"Ошибка добавления аккаунта '{account.Name}': {ex.Message}");
                    }
                }
        
                AnsiConsole.MarkupLine("[green]Банковские аккаунты успешно импортированы![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка импорта: {Markup.Escape(ex.Message)}[/]");
            }
        }

        private void Export()
        {
            try
            {
                var accounts = _accountRepository.GetAllBankAccounts().ToList();
        
                if (accounts.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Нет счетов для экспорта.[/]");
                    return;
                }
        
                var filePath = AnsiConsole.Ask<string>("Введите путь для сохранения файла:");
                var format = PromptFormatSelection();
        
                _accountImportExportFacade.Export(accounts, format, filePath);
        
                AnsiConsole.MarkupLine($"[green]Счета успешно экспортированы в файл '{Markup.Escape(filePath)}' в формате {Markup.Escape(format)}.[/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка экспорта: {Markup.Escape(ex.Message)}[/]");
            }
        }
        private void RecalculateAccount()
        {
            try
            {
                var accounts = _accountRepository.GetAllBankAccounts().ToList();

                if (accounts.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Нет счетов для пересчета.[/]");
                    return;
                }

                var accountToRecalculate = AnsiConsole.Prompt(
                    new SelectionPrompt<BankAccount>()
                        .Title("Выберите счет для пересчета баланса:")
                        .AddChoices(accounts));

                _accountRepository.RecalculateAccountBalance(accountToRecalculate.Id);

                if (accountToRecalculate.Name != null)
                    AnsiConsole.MarkupLine("[green]Счет {0} успешно пересчитан![/]",
                        Markup.Escape(accountToRecalculate.Name));
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine("[red]Ошибка пересчета счета: {0}[/]", Markup.Escape(ex.Message));
            }
        }
    }
}