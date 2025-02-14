namespace ZooLibrary.Interfaces
{
    /// <summary>
    /// Интерфейс для описания инвентаризационной единицы.
    /// </summary>
    public interface IInventory
    {
        /// <summary>
        /// Номер или количество инвентаризационной вещи.
        /// </summary>
        int Number { get; set; }
    }
}