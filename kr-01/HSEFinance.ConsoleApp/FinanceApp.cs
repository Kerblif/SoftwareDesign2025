using System.Threading.Tasks;
using HSEFinance.Lib.Application.Analytics;
using Spectre.Console;

namespace HSEFinance.ConsoleApp
{
    public class FinanceApp
    {
        private readonly AccountManagerFacade _accountFacade;
        private readonly CategoryManagerFacade _categoryFacade;
        private readonly OperationManagerFacade _operationFacade;
        private readonly AnalyticsFacade _analyticsFacade;

        public FinanceApp(OperationManagerFacade operationManager, AccountManagerFacade accountManager, CategoryManagerFacade categoryManager, AnalyticsFacade analyticsFacade)
        {
            _accountFacade = accountManager;
            _categoryFacade = categoryManager;
            _operationFacade = operationManager;
            _analyticsFacade = analyticsFacade;
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
                        .AddChoices("Управление счетами", "Управление категориями", "Управление операциями", "Аналитика", "Выйти"));

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
                    case "Аналитика":
                        _analyticsFacade.ShowMenu();
                        break;
                    case "Выйти":
                        AnsiConsole.MarkupLine("[yellow]До свидания![/]");
                        return;
                }
            }
        }
    }
}