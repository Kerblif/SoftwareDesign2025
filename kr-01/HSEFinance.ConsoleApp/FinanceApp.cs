using System.Threading.Tasks;
using Spectre.Console;

namespace HSEFinance.ConsoleApp
{
    public class FinanceApp
    {
        private readonly AccountManagerFacade _accountFacade;
        private readonly CategoryManagerFacade _categoryFacade;
        private readonly OperationManagerFacade _operationFacade;

        public FinanceApp(
            AccountManagerFacade accountFacade,
            CategoryManagerFacade categoryFacade,
            OperationManagerFacade operationFacade)
        {
            _accountFacade = accountFacade;
            _categoryFacade = categoryFacade;
            _operationFacade = operationFacade;
        }

        public void Run()
        {
            Console.Clear();
            AnsiConsole.Write(new FigletText("HSE Finance")
                .Centered()
                .Color(Color.Green));
            
            while (true)
            {
                // Вывод главного меню
                var choice = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("[green]Выберите действие:[/]")
                        .AddChoices("Управление счетами", "Управление категориями", "Управление операциями", "Выйти"));

                switch (choice)
                {
                    case "Управление счетами":
                        _accountFacade.ShowMenu();
                        break;
                    case "Управление категориями":
                        _categoryFacade.ShowMenu();
                        break;
                    case "Управление операциями":
                        _operationFacade.ShowMenu();
                        break;
                    case "Выйти":
                        AnsiConsole.MarkupLine("[yellow]До свидания![/]");
                        return;
                }
            }
        }
    }
}