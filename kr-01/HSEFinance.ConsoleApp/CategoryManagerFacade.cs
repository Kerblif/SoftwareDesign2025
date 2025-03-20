using Spectre.Console;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;

namespace HSEFinance.ConsoleApp
{
    public class CategoryManagerFacade
    {
        private readonly ICategoryRepository _categoryRepository;

        public CategoryManagerFacade(ICategoryRepository categoryRepository)
        {
            _categoryRepository = new CategoryRepositoryProxy(categoryRepository);
        }

        public void ShowMenu()
        {
            while (true)
            {
                var choice = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("[green]Управление категориями[/]")
                        .AddChoices("Добавить категорию", "Показать все категории", "Удалить категорию", "Редактировать категорию", "Назад"));

                switch (choice)
                {
                    case "Добавить категорию":
                        AddCategory();
                        break;

                    case "Показать все категории":
                        ShowAllCategories();
                        break;

                    case "Удалить категорию":
                        DeleteCategory();
                        break;
                    
                    case "Редактировать категорию":
                        EditCategory();
                        break;

                    case "Назад":
                        return; // Возврат в предыдущее меню
                }
            }
        }

        private void AddCategory()
        {
            try
            {
                // Запрос данных новой категории (тип и название)
                var type = AnsiConsole.Prompt(
                    new SelectionPrompt<string>()
                        .Title("Выберите тип категории:")
                        .AddChoices("Доход", "Расход"));

                var categoryName = AnsiConsole.Ask<string>("Введите название категории:");

                // Преобразование типа категории
                var itemType = type == "Доход" ? ItemType.Income : ItemType.Expense;

                // Создание категории в репозитории
                _categoryRepository.CreateCategory(itemType, categoryName);

                AnsiConsole.MarkupLine($"[green]Категория '{categoryName}' успешно добавлена![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка добавления категории: {ex.Message}[/]");
            }
        }

        private void ShowAllCategories()
        {
            try
            {
                // Получение всех категорий из репозитория
                var categories = _categoryRepository.GetAllCategories().ToList();

                if (categories.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Список категорий пуст.[/]");
                    return;
                }

                // Отображение данных о категориях
                AnsiConsole.MarkupLine("[green]Список категорий:[/]");

                var table = new Table()
                    .AddColumn("ID")
                    .AddColumn("Тип")
                    .AddColumn("Название");

                foreach (var category in categories)
                {
                    table.AddRow(category.Id.ToString(), category.Type.ToString(), category.Name);
                }

                AnsiConsole.Render(table);
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка отображения категорий: {ex.Message}[/]");
            }
        }

        private void DeleteCategory()
        {
            try
            {
                // Получение списка всех категорий
                var categories = _categoryRepository.GetAllCategories().ToList();

                if (categories.Count == 0)
                {
                    AnsiConsole.MarkupLine("[yellow]Нет категорий для удаления.[/]");
                    return;
                }

                // Выбор категории для удаления
                var categoryToDelete = AnsiConsole.Prompt(
                    new SelectionPrompt<Category>()
                        .Title("Выберите категорию для удаления:")
                        .AddChoices(categories));

                // Удаление категории из репозитория
                _categoryRepository.DeleteCategory(categoryToDelete.Id);

                AnsiConsole.MarkupLine($"[green]Категория '{categoryToDelete.Name}' успешно удалена![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка удаления категории: {ex.Message}[/]");
            }
        }
        
        private void EditCategory()
        {
            try
            {
                var categories = _categoryRepository.GetAllCategories().ToList();

                if (!categories.Any())
                {
                    AnsiConsole.MarkupLine("[yellow]Нет категорий для редактирования.[/]");
                    return;
                }

                var categoryToEdit = AnsiConsole.Prompt(
                    new SelectionPrompt<Category>()
                        .Title("Выберите категорию для редактирования:")
                        .AddChoices(categories));

                var newName = AnsiConsole.Ask<string>($"Введите новое название категории (текущее: {categoryToEdit.Name}):");
                var newType = AnsiConsole.Prompt(
                    new SelectionPrompt<ItemType>()
                        .Title($"Выберите новый тип категории (текущий: {categoryToEdit.Type}):")
                        .AddChoices(ItemType.Income, ItemType.Expense));

                categoryToEdit.Name = newName;
                categoryToEdit.Type = newType;

                _categoryRepository.UpdateCategory(categoryToEdit);

                AnsiConsole.MarkupLine($"[green]Категория '{categoryToEdit.Name}' успешно обновлена![/]");
            }
            catch (Exception ex)
            {
                AnsiConsole.MarkupLine($"[red]Ошибка редактирования категории: {ex.Message}[/]");
            }
        }
    }
}