using ZooLibrary.Animals;
using ZooLibrary.Things;

namespace ZooLibrary.ZooEntities
{
    /// <summary>
    /// Класс для описания зоопарка.
    /// </summary>
    public class Zoo
    {
        private static readonly Lazy<Zoo> _instance = new Lazy<Zoo>(() => new Zoo());

        private Zoo()
        {
        }
        
        public static Zoo Instance => _instance.Value;
        
        public List<Animal> Animals { get; set; } = new();
        public List<Thing> Inventory { get; set; } = new();
        public IAnimalChecker Clinic { get; set; } = new VetClinic();

        /// <summary>
        /// Попытаться добавить животное в зоопарк (через проверку клиники).
        /// </summary>
        public bool TryAddAnimal(Animal animal)
        {
            if (Clinic.Check(animal))
            {
                Animals.Add(animal);
                return true;
            }
            return false;
        }

        /// <summary>
        /// Добавить вещь в инвентарь.
        /// </summary>
        public void AddInventory(Thing thing)
        {
            Inventory.Add(thing);
        }
    }
}