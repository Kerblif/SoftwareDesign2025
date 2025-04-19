using HSEFinance.Lib.Application.Analytics;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;
using Spectre.Console;

namespace HSEFinance.ConsoleApp
{
    public class AnalyticsFacade
    {
        private readonly IOperationRepository _operationRepository;

        public AnalyticsFacade(IOperationRepository operationRepository)
        {
            _operationRepository = new OperationRepositoryProxy(operationRepository);
        }
        
        public void ShowMenu()
        {
            while (true)
            {
                var choice = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("Выберите аналитический отчет:")
                        .AddChoices(new[] {
                            "Разница доходов и расходов за период",
                            "Группировка доходов/расходов по категориям",
                            "Средняя сумма доходов/расходов",
                            "Назад"
                        }));

                switch (choice)
                {
                    case "Разница доходов и расходов за период":
                        ShowIncomeExpenseDifference();
                        break;
                    case "Группировка доходов/расходов по категориям":
                        ShowCategoryGrouping();
                        break;
                    case "Средняя сумма доходов/расходов":
                        ShowAverageStatistics();
                        break;
                    case "Назад":
                        return;
                }
            }
        }

        private void ShowIncomeExpenseDifference()
        {
            var startDate = AnsiConsole.Ask<DateTime>("Введите [green]начальную дату[/] (в формате ГГГГ-ММ-ДД):");
            var endDate = AnsiConsole.Ask<DateTime>("Введите [green]конечную дату[/] (в формате ГГГГ-ММ-ДД):");

            var visitor = new IncomeExpenseDifferenceVisitor(startDate, endDate);

            _operationRepository.Accept(visitor);

            AnsiConsole.MarkupLine("[blue]Доходы:[/] {0}", visitor.TotalIncome);
            AnsiConsole.MarkupLine("[blue]Расходы:[/] {0}", visitor.TotalExpense);
            AnsiConsole.MarkupLine("[blue]Разница:[/] {0}", visitor.CalculateDifference());
        }

        private void ShowCategoryGrouping()
        {
            var visitor = new CategoryGroupingVisitor();

            _operationRepository.Accept(visitor);

            var table = new Table()
                .AddColumn("[green]Категория[/]")
                .AddColumn("[blue]Доходы[/]")
                .AddColumn("[red]Расходы[/]");

            foreach (var income in visitor.IncomeByCategory)
            {
                var category = income.Key;
                var incomeValue = income.Value;
                var expenseValue = visitor.ExpenseByCategory.GetValueOrDefault(category);

                table.AddRow(category.ToString(), incomeValue.ToString("F2"), expenseValue.ToString("F2"));
            }

            AnsiConsole.Write(table);
        }

        private void ShowAverageStatistics()
        {
            var visitor = new AverageOperationVisitor();

            _operationRepository.Accept(visitor);

            AnsiConsole.MarkupLine("[blue]Средний доход:[/] {0:F2}", visitor.GetAverageIncome());
            AnsiConsole.MarkupLine("[red]Средний расход:[/] {0:F2}", visitor.GetAverageExpense());
        }
    }
}