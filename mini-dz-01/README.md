# Проект ZooLibrary

Проект представляет собой библиотеку для работы с концепцией зоопарка. Основные функциональные элементы включают обработку животных, их создание с помощью фабричных методов, проверку здоровья животных и управление зоопарком.

В данном README описано, где и как применяются принципы SOLID (англ. Single Responsibility, Open-Closed, Liskov Substitution, Interface Segregation, Dependency Inversion) в коде проекта.

---

## Применение принципов SOLID

### 1. **Single Responsibility Principle (Принцип единственной ответственности)**

Каждый класс в данном проекте имеет только одну ответственность.

- **Пример:**
  Класс `VetClinic` отвечает исключительно за проверку здоровья животных и никак не участвует, например, в их создании или управлении зоопарком:
  ```c#
  public class VetClinic : IAnimalChecker
  {
      public bool Check(Animal animal)
      {
          return animal.IsHealthy;
      }
  }
  ```
  Аналогично, класс `Zoo` отвечает за хранение информации про животных и инвентарь, а также за их добавление в зоопарк.

---

### 2. **Open-Closed Principle (Принцип открытости/закрытости)**

Классы открыты для расширения, но закрыты для модификации. Новая функциональность может быть добавлена через наследование и реализацию интерфейсов, не изменяя уже существующий код.

- **Пример:**
  Базовый абстрактный класс `AnimalFactoryBase<TAnimal>` реализует общую логику для всех фабрик по созданию животных. Для добавления нового вида животного (например, `Giraffe`) нужно лишь создать новый класс-фабрику, унаследованный от `AnimalFactoryBase<TAnimal>`, без необходимости модификации существующих классов:
  ```c#
  public abstract class AnimalFactoryBase<TAnimal> : IAutoRegisteredAnimalFactory where TAnimal : Animal
  {
      public Type GetAnimalType() => typeof(TAnimal);
      public abstract Animal CreateAnimal(int food, int number, bool isHealthy);
  }
  ```

- Конкретные фабрики, такие как `MonkeyFactory` и `RabbitFactory`, расширяют функциональность, не изменяя базовые классы.

---

### 3. **Liskov Substitution Principle (Принцип подстановки Барбары Лисков)**

Классы-наследники могут заменять своих базовых предков, не нарушая работоспособность программы.

- **Пример:**
  Все фабрики, унаследованные от `AnimalFactoryBase<TAnimal>`, взаимозаменяемы, поскольку принимаются и используются через их общий интерфейс `IAutoRegisteredAnimalFactory`:
  ```c#
  public interface IAutoRegisteredAnimalFactory
  {
      Type GetAnimalType();
      Animal CreateAnimal(int food, int number, bool isHealthy);
  }
  ```

- Аналогично, любой объект, реализующий `IAnimalChecker` (например, `VetClinic`), может быть использован для проверки животных в классе `Zoo`:
  ```c#
  public IAnimalChecker Clinic { get; set; } = new VetClinic();
  ```

---

### 4. **Interface Segregation Principle (Принцип разделения интерфейсов)**

Интерфейсы в проекте небольшие и специализированные, чтобы классы должны были реализовывать только те методы, которые им действительно нужны.

- **Пример:**
  Интерфейс `IAutoRegisteredAnimalFactory` содержит только два метода, связанных с функциональностью фабрики:
  ```c#
  public interface IAutoRegisteredAnimalFactory
  {
      Type GetAnimalType();
      Animal CreateAnimal(int food, int number, bool isHealthy);
  }
  ```
  Классы, реализующие этот интерфейс (например, `MonkeyFactory`, `RabbitFactory`), не обязаны определять что-либо лишнее.

---

### 5. **Dependency Inversion Principle (Принцип инверсии зависимостей)**

Зависимости определены через абстракции, а не конкретные реализации.

- **Пример:**
  В классе `Zoo` для проверки животных используется интерфейс `IAnimalChecker`, а не конкретный класс:
  ```c#
  public IAnimalChecker Clinic { get; set; } = new VetClinic();
  ```
  Это позволяет легко подменять реализацию. Например, вместо `VetClinic` можно использовать другой класс, реализующий `IAnimalChecker`.

- DI-контейнер в классе `AppServices` регистрирует зависимости через внедрение:
  ```c#
  services.AddSingleton<IAnimalChecker, VetClinic>();
  ```

---

## Тестирование

Принципы проектирования также применялись для обеспечения тестируемости кода. Например, классы с внедрением интерфейсов легко тестировать через мок-объекты. Тесты для `VetClinic` показывают, что поведение класса однозначно и соответствует ожидаемому:

```c#
[Test]
public void VetClinic_ShouldReturnTrue_ForHealthyAnimal()
{
    var animal = new Monkey(5, 1, true, 3);
    var clinic = new VetClinic();
    var result = clinic.Check(animal);
    Assert.IsTrue(result);
}
```

---

## Выводы

Применение принципов SOLID в данном проекте позволило достичь модульности, читаемости и расширяемости кода. Каждый компонент выполняет строго определенную задачу, новая функциональность добавляется без модификации существующего кода, а интерфейсы делают систему гибкой.