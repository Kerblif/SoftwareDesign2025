using Spectre.Console;
using ZooLibrary.Animals;
using ZooLibrary.ZooEntities;
using ZooLibrary.Things;

using SpectreTable = Spectre.Console.Table;
using ZooTable = ZooLibrary.Things.Table;

namespace ZooConsoleApp
{
    public class Program
    {
        private static Zoo _zoo = new();

        public static void Main(string[] args)
        {
            AnsiConsole.Write(
                new FigletText("Moscow Zoo")
                    .Color(Color.Green));

            while (true)
            {
                var choice = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("[yellow]Выберите действие:[/]")
                        .AddChoices("Добавить животное", "Добавить вещь", "Показать отчет", "Список контактного зоопарка", "Показать инвентарь", "Выход"));

                switch (choice)
                {
                    case "Добавить животное":
                        AddAnimal();
                        break;
                    case "Добавить вещь":
                        AddThing();
                        break;
                    case "Показать отчет":
                        ShowAnimalReport();
                        break;
                    case "Список контактного зоопарка":
                        ShowInteractiveAnimals();
                        break;
                    case "Показать инвентарь":
                        ShowInventory();
                        break;
                    case "Выход":
                        return;
                }

                AnsiConsole.MarkupLine("[green]Нажмите Enter, чтобы продолжить...[/]");
                Console.ReadLine();
                Console.Clear();
            }
        }

        private static void AddAnimal()
        {
            var species = AnsiConsole.Prompt(
                new SelectionPrompt<string>()
                    .Title("Выберите [green]вид животного[/]:")
                    .AddChoices("Обезьяна", "Кролик"));

            var food = AnsiConsole.Ask<int>("Сколько [green]килограммов еды[/] нужно этому животному в день?");
            var number = AnsiConsole.Ask<int>("Введите [green]инвентарный номер[/]:");
            var isHealthy = AnsiConsole.Confirm("Это животное [green]здорово[/]?");

            Animal animal;
            try
            {
                animal = species switch
                {
                    "Обезьяна" => new Monkey(
                        food,
                        number,
                        isHealthy,
                        AnsiConsole.Ask<int>("Введите [green]уровень доброты[/] (0-10):")),
                    "Кролик" => new Rabbit(
                        food,
                        number,
                        isHealthy,
                        AnsiConsole.Ask<int>("Введите [green]уровень доброты[/] (0-10):")),
                    _ => throw new InvalidOperationException()
                };
            }
            catch (InvalidOperationException)
            {
                AnsiConsole.MarkupLine("[red]Ошибка: произошла недопустимая операция.[/]");
                return;
            }
            catch (ArgumentException)
            {
                AnsiConsole.MarkupLine("[red]Ошибка: недопустимый аргумент. Проверьте введенные данные.[/]");
                return;
            }

            if (_zoo.TryAddAnimal(animal))
            {
                AnsiConsole.MarkupLine($"[green]Животное добавлено успешно![/]");
            }
            else
            {
                AnsiConsole.MarkupLine($"[red]Животное не прошло проверку здоровья.[/]");
            }
        }

        private static void AddThing()
        {
            var thingType = AnsiConsole.Prompt(
                new SelectionPrompt<string>()
                    .Title("Выберите [green]тип вещи[/]:")
                    .AddChoices("Стол", "Компьютер"));

            var number = AnsiConsole.Ask<int>("Введите [green]инвентарный номер[/]:");

            Thing thing;
            try
            {
                thing = thingType switch
                {
                    "Стол" => new ZooTable(number),
                    "Компьютер" => new Computer(number),
                    _ => throw new InvalidOperationException("Не поддерживаемый тип вещи.")
                };
            }
            catch (InvalidOperationException)
            {
                AnsiConsole.MarkupLine("[red]Ошибка: не поддерживаемый тип вещи.[/]");
                return;
            }

            _zoo.AddInventory(thing);
            AnsiConsole.MarkupLine("[green]Вещь добавлена успешно![/]");
        }

        private static void ShowAnimalReport()
        {
            if (!_zoo.Animals.Any())
            {
                AnsiConsole.MarkupLine("[yellow]Животные отсутствуют на балансе![/]");
                return;
            }

            var table = new SpectreTable();
            table.AddColumn("Инвентарный номер");
            table.AddColumn("Вид");
            table.AddColumn("Потребление пищи (кг/день)");

            foreach (var animal in _zoo.Animals)
            {
                table.AddRow(animal.Number.ToString(), animal.GetType().Name, animal.Food.ToString());
            }

            AnsiConsole.Write(table);

            var totalFood = _zoo.Animals.Sum(a => a.Food);
            AnsiConsole.MarkupLine($"\n[green]Всего животных:[/] {_zoo.Animals.Count}");
            AnsiConsole.MarkupLine($"[green]Общее потребление пищи в день (кг):[/] {totalFood}");
        }

        private static void ShowInteractiveAnimals()
        {
            var interactiveAnimals = _zoo.Animals.OfType<Herbo>().Where(h => h.IsInteractive()).ToList();

            if (!interactiveAnimals.Any())
            {
                AnsiConsole.MarkupLine("[yellow]Нет животных, подходящих для контактного зоопарка![/]");
                return;
            }

            var table = new SpectreTable();
            table.AddColumn("Инвентарный номер");
            table.AddColumn("Вид");

            foreach (var animal in interactiveAnimals)
            {
                table.AddRow(animal.Number.ToString(), animal.GetType().Name);
            }

            AnsiConsole.Write(table);
        }

        private static void ShowInventory()
        {
            if (!_zoo.Inventory.Any())
            {
                AnsiConsole.MarkupLine("[yellow]Инвентарь пуст![/]");
                return;
            }

            var table = new SpectreTable();
            table.AddColumn("Инвентарный номер");
            table.AddColumn("Наименование");

            foreach (var item in _zoo.Inventory)
            {
                table.AddRow(item.Number.ToString(), item.GetType().Name);
            }

            AnsiConsole.Write(table);
        }
    }
}