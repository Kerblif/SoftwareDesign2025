using HSEFinance.Lib.Application.Commands;
using Spectre.Console;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;

namespace HSEFinance.ConsoleApp
{
    public class OperationManagerFacade
    {
        private readonly IOperationRepository _operationRepository;

        public OperationManagerFacade(IOperationRepository operationRepository)
        {
            _operationRepository = new OperationRepositoryProxy(operationRepository);
        }

        public void ShowMenu()
        {
            while (true)
            {
                // Главное меню управления операциями
                var choice = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("[green]Управление операциями[/]")
                        .AddChoices(
                            "Добавить операцию",
                            "Показать все операции",
                            "Удалить операцию",
                            "Редактировать операцию",
                            "Назад"));

                switch (choice)
                {
                    case "Добавить операцию":
                        AddOperation();
                        break;

                    case "Показать все операции":
                        ShowAllOperations();
                        break;

                    case "Удалить операцию":
                        DeleteOperation();
                        break;
                    
                    case "Редактировать операцию":
                        EditOperation();
                        break;


                    case "Назад":
                        return;
                }
            }
        }

        private void AddOperation()
        {
            try
            {
                // Выбор типа операции
                var type = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("Выберите тип операции:")
                        .AddChoices("Доход", "Расход"));

                var bankAccountId = AnsiConsole.Ask<Guid>("Введите ID банковского счета:");
                var amount = AnsiConsole.Ask<decimal>("Введите сумму операции:");
                var date = AnsiConsole.Ask<DateTime>("Введите дату операции (в формате YYYY-MM-DD):");
                var categoryId = AnsiConsole.Ask<Guid>("Введите ID категории:");
                var description = AnsiConsole.Ask<string>("Введите описание операции (необязательно):", "");

                var itemType = type == "Доход" ? ItemType.Income : ItemType.Expense;

                ICommand command = new CreateOperationCommand(_operationRepository, itemType, bankAccountId, amount, date, categoryId, description);
                command = new TimedCommand(command, time => Console.WriteLine("Выполнено за {0:D} милисекунд", time.Milliseconds));
                command.Execute();
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка добавления операции: {ex.Message}[/]");
            }
        }

        private void ShowAllOperations()
        {
            try
            {
                var operations = _operationRepository.GetAllOperations().ToList();

                if (operations.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Список операций пуст.[/]");
                    return;
                }

                // Отображение операций в табличной форме
                AnsiConsole.MarkupLine("[green]Список операций:[/]");

                var table = new Table()
                    .AddColumn("ID")
                    .AddColumn("Тип")
                    .AddColumn("Счет")
                    .AddColumn("Сумма")
                    .AddColumn("Дата")
                    .AddColumn("Категория")
                    .AddColumn("Описание");

                foreach (var operation in operations)
                {
                    table.AddRow(
                        operation.Id.ToString(),
                        operation.Type.ToString(),
                        operation.BankAccountId.ToString(),
                        operation.Amount.ToString("C"),
                        operation.Date.ToString("yyyy-MM-dd"),
                        operation.CategoryId.ToString(),
                        operation.Description ?? "—");
                }

                AnsiConsole.Render(table);
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка отображения операций: {ex.Message}[/]");
            }
        }

        private void DeleteOperation()
        {
            try
            {
                var operations = _operationRepository.GetAllOperations().ToList();

                if (operations.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Нет операций для удаления.[/]");
                    return;
                }

                // Выбор операции для удаления
                var operationToDelete = AnsiConsole.Prompt(
                    new SelectionPrompt<Operation>()
                        .Title("Выберите операцию для удаления:")
                        .AddChoices(operations));

                // Создание команды удаления операции
                ICommand command = new DeleteOperationCommand(_operationRepository, operationToDelete.Id);
                command = new TimedCommand(command, time => Console.WriteLine("Выполнено за {0:D} милисекунд", time.Milliseconds));
                command.Execute();

                AnsiConsole.MarkupLine("[green]Операция на удаление добавлена в список команд![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка удаления операции: {ex.Message}[/]");
            }
        }
        
        private void EditOperation()
        {
            try
            {
                var operations = _operationRepository.GetAllOperations().ToList();

                if (!operations.Any())
                {
                    AnsiConsole.MarkupLine("[yellow]Нет операций для редактирования.[/]");
                    return;
                }

                var operationToEdit = AnsiConsole.Prompt(
                    new SelectionPrompt<Operation>()
                        .Title("Выберите операцию для редактирования:")
                        .AddChoices(operations));
                
                operationToEdit.CategoryId = AnsiConsole.Ask<Guid>($"Введите новый ID категории (текущий: {operationToEdit.CategoryId}):", operationToEdit.CategoryId);
                operationToEdit.Description = AnsiConsole.Ask<string>($"Введите новое описание операции (текущее: {operationToEdit.Description ?? "N/A"}):", operationToEdit.Description);

                _operationRepository.UpdateOperation(operationToEdit);

                AnsiConsole.MarkupLine($"[green]Операция успешно обновлена![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка редактирования операции: {ex.Message}[/]");
            }
        }
    }
}
