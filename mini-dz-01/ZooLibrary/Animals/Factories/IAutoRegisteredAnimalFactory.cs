namespace ZooLibrary.Animals.Factories
{
    public interface IAutoRegisteredAnimalFactory
    {
        /// <summary>
        /// Возвращает тип животного, который создает эта фабрика.
        /// </summary>
        Type GetAnimalType();

        /// <summary>
        /// Создание животного.
        /// </summary>
        Animal CreateAnimal(int food, int number, bool isHealthy);
    }
}