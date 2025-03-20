using Spectre.Console;
using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;

namespace HSEFinance.ConsoleApp
{
    public class AccountManagerFacade
    {
        private readonly IBankAccountFactory _accountFactory;
        private readonly IAccountRepository _accountRepository;

        public AccountManagerFacade(IBankAccountFactory accountFactory, IAccountRepository accountRepository)
        {
            _accountFactory = accountFactory;
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
                        .AddChoices("Добавить счет", "Показать все счета", "Удалить счет", "Редактировать счет", "Назад"));

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

                AnsiConsole.MarkupLine($"[green]Счет '{name}' успешно добавлен![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка добавления счета: {ex.Message}[/]");
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
                AnsiConsole.MarkupLine($"[red]Ошибка отображения счетов: {ex.Message}[/]");
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

                AnsiConsole.MarkupLine($"[green]Счет '{accountToDelete.Name}' успешно удален![/]");
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

                var newName = AnsiConsole.Ask<string>($"Введите новое название для счета (текущее: {accountToEdit.Name}):");
                
                accountToEdit.Name = newName;

                _accountRepository.UpdateBankAccount(accountToEdit);

                AnsiConsole.MarkupLine($"[green]Счет '{accountToEdit.Name}' успешно обновлен![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка редактирования счета: {ex.Message}[/]");
            }
        }
    }
}