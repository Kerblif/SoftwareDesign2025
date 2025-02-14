using ZooLibrary.Interfaces;

namespace ZooLibrary.Animals
{
    /// <summary>
    /// Базовый класс, описывающий животное.
    /// </summary>
    public abstract class Animal : IAlive, IInventory
    {
        private int _food;
        
        /// <inheritdoc />
        public int Food
        {
            get => _food;
            set
            {
                if (value < 0)
                {
                    throw new ArgumentException("Food cannot be negative.");
                }
                
                _food = value;
            }
        }
        
        /// <inheritdoc />
        public int Number { get; set; }

        /// <summary>
        /// Здоровье животного (проверяется ветеринарной клиникой).
        /// </summary>
        public bool IsHealthy { get; set; }

        protected Animal(int food, int number, bool isHealthy)
        {
            Food = food;
            Number = number;
            IsHealthy = isHealthy;
        }
        
        public override string ToString()
        {
            return $"{GetType().Name}: {Number}";
        }
    }
}