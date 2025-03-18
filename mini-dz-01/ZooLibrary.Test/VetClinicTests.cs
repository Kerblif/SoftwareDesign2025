using NUnit.Framework;
using ZooLibrary.Animals;
using ZooLibrary.ZooEntities;

namespace ZooLibrary.Test;

[TestFixture]
public class VetClinicTests
{
    [Test]
    public void VetClinic_ShouldReturnTrue_ForHealthyAnimal()
    {
        // Arrange
        var animal = new Monkey(5, 1, true, 3);
        var clinic = new VetClinic();
        
        // Act
        var result = clinic.Check(animal);
        
        // Assert
        Assert.IsTrue(result);
    }
    
    [Test]
    public void VetClinic_ShouldReturnFalse_ForUnhealthyAnimal()
    {
        // Arrange
        var animal = new Rabbit(8, 1, false, 7);
        var clinic = new VetClinic();
        
        // Act
        var result = clinic.Check(animal);
        
        // Assert
        Assert.IsFalse(result);
    }
}