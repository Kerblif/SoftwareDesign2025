namespace ZooLibrary.Animals.Factories
{
    public class WolfFactory : AnimalFactoryBase<Wolf>
    {
        public override Animal CreateAnimal(int food, int number, bool isHealthy)
        {
            return new Wolf(food, number, isHealthy);
        }
    }
}