namespace ZooLibrary.Animals
{
    /// <summary>
    /// Класс, описывающий травоядных животных.
    /// </summary>
    public abstract class Herbo : Animal
    {
        private const int MinimumKindness = 0;
        private const int MaximumKindness = 10;

        private int _kindness;

        /// <summary>
        /// Уровень доброты (от 0 до 10).
        /// </summary>
        public int Kindness
        {
            get => _kindness;
            set
            {
                if (value < MinimumKindness || value > MaximumKindness)
                {
                    throw new ArgumentException($"Kindness must be between {MinimumKindness} and {MaximumKindness}");
                }
                _kindness = value;
            }
        }

        protected Herbo(int food, int number, bool isHealthy, int kindness)
            : base(food, number, isHealthy)
        {
            Kindness = kindness;
        }

        /// <summary>
        /// Проверяет, доступно ли животное для интерактива.
        /// </summary>
        /// <returns></returns>
        public bool IsInteractive() => Kindness > 5;
    }
}