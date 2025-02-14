

// Атрибут для указания локализованного имени
namespace ZooLibrary.Animals;

[AttributeUsage(AttributeTargets.Class, Inherited = false)]
public class AnimalNameAttribute : Attribute
{
    public string Name { get; }
    public AnimalNameAttribute(string name) => Name = name;
}