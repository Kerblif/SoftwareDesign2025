using ZooLibrary.Animals;

namespace ZooLibrary.ZooEntities
{
    /// <summary>
    /// Класс, описывающий ветеринарную клинику.
    /// </summary>
    public class VetClinic
    {
        public bool CheckAnimal(Animal animal)
        {
            // Проверка здоровья. Написан простой пример.
            return animal.IsHealthy;
        }
    }
}