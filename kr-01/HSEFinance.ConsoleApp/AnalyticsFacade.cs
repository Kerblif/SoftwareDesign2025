using HSEFinance.Lib.Application.Analytics;
using HSEFinance.Lib.Domain.Repositories;
using Spectre.Console;

namespace HSEFinance.ConsoleApp
{
    public class AnalyticsFacade
    {
        private readonly IOperationRepository _operationRepository;

        public AnalyticsFacade(IOperationRepository operationRepository)
        {
            _operationRepository = operationRepository;
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

            foreach (var operation in _operationRepository.GetAllOperations())
            {
                operation.Accept(visitor);
            }

            AnsiConsole.MarkupLine($"[blue]Доходы:[/] {visitor.TotalIncome}");
            AnsiConsole.MarkupLine($"[blue]Расходы:[/] {visitor.TotalExpense}");
            AnsiConsole.MarkupLine($"[blue]Разница:[/] {visitor.CalculateDifference()}");
        }

        private void ShowCategoryGrouping()
        {
            var visitor = new CategoryGroupingVisitor();

            foreach (var operation in _operationRepository.GetAllOperations())
            {
                operation.Accept(visitor);
            }

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

            foreach (var operation in _operationRepository.GetAllOperations())
            {
                operation.Accept(visitor);
            }

            AnsiConsole.MarkupLine($"[blue]Средний доход:[/] {visitor.GetAverageIncome():F2}");
            AnsiConsole.MarkupLine($"[red]Средний расход:[/] {visitor.GetAverageExpense():F2}");
        }
    }
}