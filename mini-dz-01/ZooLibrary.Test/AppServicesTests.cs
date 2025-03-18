using NUnit.Framework;
using ZooLibrary.Services;

namespace ZooLibrary.Test;

[TestFixture]
public class AppServicesTests
{
    [Test]
    public void Services_ShouldNotBeNull()
    {
        // Act
        var serviceProvider = AppServices.Services;
        
        // Assert
        Assert.IsNotNull(serviceProvider);
    }
}