namespace ZooLibrary.Animals
{ 
    /// <summary>
    /// Класс для описания тигра.
    /// </summary>
    public class Tiger : Predator
    {
        public Tiger(int food, int number, bool isHealthy)
            : base(food, number, isHealthy)
        {
        }
    }
}