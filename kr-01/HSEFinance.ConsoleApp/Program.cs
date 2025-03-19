using Spectre.Console;

class Program
{
    static void Main()
    {
        // Приветствие
        AnsiConsole.Write(new FigletText("HSE Finance")
           .Centered()
           .Color(Color.Green));

        // Основное меню
        while (true)
        {
            var choice = AnsiConsole.Prompt(
                new SelectionPrompt<string>()
                    .Title("Выберите действие:")
                    .AddChoices(new[]
                    {
                        "1. Просмотр счетов",
                        "2. Добавление операций",
                        "3. Аналитика",
                        "4. Управление категориями",
                        "5. Выход"
                    }));

            switch (choice)
            {
                case "1. Просмотр счетов":
                    ShowAccounts();
                    break;
                case "2. Добавление операций":
                    AddOperation();
                    break;
                case "3. Аналитика":
                    ShowAnalytics();
                    break;
                case "4. Управление категориями":
                    ManageCategories();
                    break;
                case "5. Выход":
                    AnsiConsole.Markup("[bold green]До скорой встречи![/]");
                    return;
                default:
                    AnsiConsole.Markup("[red]Неверный выбор![/]");
                    break;
            }
        }
    }

    static void ShowAccounts()
    {
        AnsiConsole.MarkupLine("[yellow]Просмотр счетов пока не реализован[/]");
    }

    static void AddOperation()
    {
        // Пример добавления операции
        AnsiConsole.MarkupLine("[yellow]> Добавление новой операции[/]");

        var type = AnsiConsole.Prompt(
            new SelectionPrompt<string>()
                .Title("Тип операции:")
                .AddChoices("Доход", "Расход"));

        var amount = AnsiConsole.Ask<decimal>("Введите сумму операции:");
        var description = AnsiConsole.Ask<string>("Введите описание операции (необязательно):");

        // Объект операции
        AnsiConsole.MarkupLine($"[green]Операция добавлена:[/] {type} на сумму {amount:C}. Описание: {description}");
    }

    static void ShowAnalytics()
    {
        AnsiConsole.MarkupLine("[blue]Аналитика пока не реализована[/]");
    }

    static void ManageCategories()
    {
        AnsiConsole.MarkupLine("[yellow]Управление категориями пока не реализовано[/]");
    }
}