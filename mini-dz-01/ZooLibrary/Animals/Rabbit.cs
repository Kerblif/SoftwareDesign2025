namespace ZooLibrary.Animals
{
    /// <summary>
    /// Класс для описания кролика.
    /// </summary>
    [AnimalName("Кролик")]
    public class Rabbit : Herbo
    {
        public Rabbit(int food, int number, bool isHealthy, int kindness)
            : base(food, number, isHealthy, kindness)
        {
        }
    }
}