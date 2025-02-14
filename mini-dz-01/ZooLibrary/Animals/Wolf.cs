namespace ZooLibrary.Animals
{
    /// <summary>
    /// Класс, описывающий волка.
    /// </summary>
    [AnimalName("Волк")]
    public class Wolf : Predator
    {
        /// <summary>
        /// Создает объект волка.
        /// </summary>
        /// <param name="food">Количество еды, потребляемое в сутки.</param>
        /// <param name="number">Номер инвентаризационной записи.</param>
        /// <param name="isHealthy">Признак здоровья волка.</param>
        public Wolf(int food, int number, bool isHealthy)
            : base(food, number, isHealthy)
        {
        }
    }
}