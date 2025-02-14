namespace ZooLibrary.Animals
{
    /// <summary>
    /// Класс, описывающий хищников.
    /// </summary>
    public abstract class Predator : Animal
    {
        protected Predator(int food, int number, bool isHealthy)
            : base(food, number, isHealthy)
        {
        }
    }
}