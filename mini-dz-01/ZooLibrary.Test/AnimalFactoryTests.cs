using NUnit.Framework;
using ZooLibrary.Animals;
using ZooLibrary.Animals.Factories;

namespace ZooLibrary.Test;

[TestFixture]
public class AnimalFactoryTests
{
    [Test]
    public void MonkeyFactory_ShouldCreateHealthyMonkey()
    {
        // Arrange
        var factory = new MonkeyFactory();
        
        // Act
        var monkey = factory.CreateAnimal(10, 2, true);
        
        // Assert
        Assert.IsInstanceOf<Monkey>(monkey);
        Assert.AreEqual(10, monkey.Food); // Проверяем корм
        Assert.AreEqual(2, monkey.Number); // Проверяем номер
        Assert.IsTrue(monkey.IsHealthy); // Проверяем здоровье
    }
    
    [Test]
    public void RabbitFactory_ShouldCreateUnhealthyRabbit()
    {
        // Arrange
        var factory = new RabbitFactory();
        
        // Act
        var rabbit = factory.CreateAnimal(15, 1, false);
        
        // Assert
        Assert.IsInstanceOf<Rabbit>(rabbit);
        Assert.AreEqual(15, rabbit.Food); // Проверка корма
        Assert.AreEqual(1, rabbit.Number); // Проверка номера
        Assert.IsFalse(rabbit.IsHealthy); // Проверка здоровья
    }
}