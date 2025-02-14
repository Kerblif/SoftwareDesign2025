namespace ZooLibrary.Interfaces
{
    /// <summary>
    /// Интерфейс, представляющий живое существо.
    /// </summary>
    public interface IAlive
    {
        /// <summary>
        /// Количество еды, потребляемой в килограммах в сутки.
        /// </summary>
        int Food { get; set; }
    }
}
