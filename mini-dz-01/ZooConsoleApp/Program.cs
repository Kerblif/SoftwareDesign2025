using Microsoft.Extensions.DependencyInjection;
using Spectre.Console;
using ZooLibrary.Animals;
using ZooLibrary.Animals.Factories;
using ZooLibrary.Services;
using ZooLibrary.ZooEntities;
using ZooLibrary.Things;

using SpectreTable = Spectre.Console.Table;
using ZooTable = ZooLibrary.Things.Table;

namespace ZooConsoleApp
{
    public class Program
    {
        private static Zoo _zoo = null!;

        public static void Main()
        {
            _zoo = AppServices.Services.GetRequiredService<Zoo>();
            
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
            var animalsRegistry = AppServices.Services.GetRequiredService<AnimalNameRegistry>();
            var factoryRegistry = AppServices.Services.GetRequiredService<AnimalFactoryRegistry>();

            var species = animalsRegistry.GetAnimalNames().ToList();

            if (!species.Any())
            {
                AnsiConsole.MarkupLine("[red]Ошибка: ни одна фабрика животных не зарегистрирована![/]");
                return;
            }
            
            var selectedSpecies= AnsiConsole.Prompt(
                new SelectionPrompt<string>()
                    .Title("Выберите [green]вид животного[/]:")
                    .AddChoices(species)); // Получаем список зарегистрированных видов

            var food = AnsiConsole.Ask<int>("Сколько [green]килограммов еды[/] нужно этому животному в день?");
            var number = AnsiConsole.Ask<int>("Введите [green]инвентарный номер[/]:");
            var isHealthy = AnsiConsole.Confirm("Это животное [green]здорово[/]?");

            var factory = factoryRegistry.GetFactory(selectedSpecies);
            if (factory == null)
            {
                AnsiConsole.MarkupLine("[red]Ошибка: фабрика для данного вида животных не найдена.[/]");
                return;
            }

            Animal animal;

            try
            {
                animal = factory.CreateAnimal(food, number, isHealthy);

                if (animal is Herbo herbo)
                {
                    var kindness = AnsiConsole.Prompt(
                        new TextPrompt<int>("Введите [green]уровень доброты[/] (0-10):")
                            .Validate((n) => n switch
                            {
                                < 0 => ValidationResult.Error("Слишком низкий уровень доброты"),
                                > 50 => ValidationResult.Error("Слишком высокий уровень доброты"),
                                _ => ValidationResult.Success()
                            })
                            .DefaultValue(herbo.Kindness));
                    herbo.Kindness = kindness;
                }
            }
            catch (ArgumentException ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка создания животного: {ex.Message}[/]");
                return;
            }

            if (_zoo.TryAddAnimal(animal))
            {
                AnsiConsole.MarkupLine($"[green]Животное успешно добавлено![/]");
            }
            else
            {
                AnsiConsole.MarkupLine($"[red]Животное не прошло проверку здоровья и не было добавлено.[/]");
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