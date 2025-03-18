namespace ZooLibrary.Animals
{
    /// <summary>
    /// Класс, описывающий обезьяну.
    /// </summary>
    [AnimalName("Обезьяна")]
    public class Monkey : Herbo
    {
        /// <summary>
        /// Создает объект обезьяны.
        /// </summary>
        /// <param name="food">Количество еды, потребляемое в сутки.</param>
        /// <param name="number">Номер инвентаризационной записи.</param>
        /// <param name="isHealthy">Признак здоровья обезьяны.</param>
        public Monkey(int food, int number, bool isHealthy, int kindness)
            : base(food, number, isHealthy, kindness)
        {
        }
    }
}