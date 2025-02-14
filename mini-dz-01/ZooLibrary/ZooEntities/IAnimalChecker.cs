using ZooLibrary.Animals;

namespace ZooLibrary.ZooEntities
{
    /// <summary>
    /// Интерфейс для проверки животных.
    /// </summary>
    public interface IAnimalChecker
    {
        bool Check(Animal animal);
    }
}