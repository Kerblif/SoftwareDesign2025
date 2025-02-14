using NUnit.Framework;
using ZooLibrary.Animals;
using ZooLibrary.Things;
using ZooLibrary.ZooEntities;
using System.Collections.Generic;

namespace ZooLibrary.Test
{
    [TestFixture]
    public class ZooArraySetTests
    {
        [SetUp]
        public void Setup()
        {
            // Сбрасываем состояние Singleton перед каждым тестом
            Zoo.Instance.Animals.Clear();
            Zoo.Instance.Inventory.Clear();
        }

        [Test]
        public void SetAnimals_ShouldUpdateAnimalList()
        {
            // Arrange
            var zoo = Zoo.Instance;
            var newAnimals = new List<Animal>
            {
                new Monkey(10, 1, true, 3),
                new Rabbit(5, 2, true, 7)
            };

            // Act
            zoo.Animals = newAnimals;

            // Assert
            Assert.AreEqual(2, zoo.Animals.Count);
            Assert.Contains(newAnimals[0], zoo.Animals);
            Assert.Contains(newAnimals[1], zoo.Animals);
        }

        [Test]
        public void SetInventory_ShouldUpdateInventoryList()
        {
            // Arrange
            var zoo = Zoo.Instance;
            var newInventory = new List<Thing>
            {
                new Computer(24),
                new Table(38)
            };

            // Act
            zoo.Inventory = newInventory;

            // Assert
            Assert.AreEqual(2, zoo.Inventory.Count);
            Assert.Contains(newInventory[0], zoo.Inventory);
            Assert.Contains(newInventory[1], zoo.Inventory);
        }

        [Test]
        public void SetEmptyAnimals_ShouldClearAnimalList()
        {
            // Arrange
            var zoo = Zoo.Instance;
            var newAnimals = new List<Animal>(); // Пустой список

            // Add initial animals to ensure clearing works
            zoo.Animals.Add(new Monkey(10, 1, true, 3));

            // Act
            zoo.Animals = newAnimals;

            // Assert
            Assert.IsEmpty(zoo.Animals);
        }

        [Test]
        public void SetEmptyInventory_ShouldClearInventoryList()
        {
            // Arrange
            var zoo = Zoo.Instance;
            var newInventory = new List<Thing>(); // Пустой список

            // Add initial items to ensure clearing works
            zoo.Inventory.Add(new Table(47));

            // Act
            zoo.Inventory = newInventory;

            // Assert
            Assert.IsEmpty(zoo.Inventory);
        }

        [Test]
        public void SetNewAnimals_ShouldReplaceExistingAnimals()
        {
            // Arrange
            var zoo = Zoo.Instance;

            var initialAnimals = new List<Animal>
            {
                new Monkey(10, 1, true, 3)
            };

            var newAnimals = new List<Animal>
            {
                new Rabbit(15, 2, true, 7)
            };

            // Set initial animals
            zoo.Animals = initialAnimals;

            // Act
            zoo.Animals = newAnimals;

            // Assert
            Assert.AreEqual(1, zoo.Animals.Count);
            Assert.AreEqual(newAnimals[0], zoo.Animals[0]); // Проверяем, что только новый список
        }

        [Test]
        public void SetNewInventory_ShouldReplaceExistingInventory()
        {
            // Arrange
            var zoo = Zoo.Instance;

            var initialInventory = new List<Thing>
            {
                new Computer(23)
            };

            var newInventory = new List<Thing>
            {
                new Table(47)
            };

            // Set initial inventory
            zoo.Inventory = initialInventory;

            // Act
            zoo.Inventory = newInventory;

            // Assert
            Assert.AreEqual(1, zoo.Inventory.Count);
            Assert.AreEqual(newInventory[0], zoo.Inventory[0]); // Проверяем, что только новый список
            Assert.AreEqual(newInventory[0].Number, zoo.Inventory[0].Number);
            Assert.AreEqual(zoo.Inventory[0].ToString(), zoo.Inventory[0].ToString());
        }
    }
}