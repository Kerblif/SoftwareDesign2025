using ZooLibrary.Animals;
using ZooLibrary.Things;

namespace ZooLibrary.ZooEntities
{
    /// <summary>
    /// Класс для описания зоопарка.
    /// </summary>
    public class Zoo
    {
        public List<Animal> Animals { get; set; } = new();
        public List<Thing> Inventory { get; set; } = new();
        public VetClinic Clinic { get; set; } = new();

        /// <summary>
        /// Попытаться добавить животное в зоопарк (через проверку клиники).
        /// </summary>
        public bool TryAddAnimal(Animal animal)
        {
            if (Clinic.CheckAnimal(animal))
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