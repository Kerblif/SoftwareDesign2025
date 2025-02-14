using ZooLibrary.Interfaces;

namespace ZooLibrary.Things
{
    /// <summary>
    /// Базовый класс для вещей.
    /// </summary>
    public abstract class Thing : IInventory
    {
        public int Number { get; set; }

        protected Thing(int number)
        {
            Number = number;
        }

        public override string ToString()
        {
            return $"{GetType().Name}: {Number}";
        }
    }
}