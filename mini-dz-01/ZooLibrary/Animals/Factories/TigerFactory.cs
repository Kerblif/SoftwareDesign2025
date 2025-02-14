namespace ZooLibrary.Animals.Factories
{
    public class TigerFactory : AnimalFactoryBase<Tiger>
    {
        public override Animal CreateAnimal(int food, int number, bool isHealthy)
        {
            return new Tiger(food, number, isHealthy);
        }
    }
}