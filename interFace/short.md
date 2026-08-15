1. Abstract Class

An abstract class is a class that can contain both:

    abstract methods → methods without implementation
    concrete methods → methods with implementation
    variables/state
    constructors

2. Interface

    An interface defines a contract/capability that a class agrees to provide.


| Feature              | Abstract Class                        | Interface                                     |
| -------------------- | ------------------------------------- | --------------------------------------------- |
| Keyword              | `abstract class`                      | `interface`                                   |
| Inheritance          | `extends`                             | `implements`                                  |
| Multiple inheritance | A class can extend only **one** class | A class can implement **multiple** interfaces |
| Constructor          | **Yes**                               | **No constructor**                            |
| Instance variables   | **Yes**                               | Normally constants (`public static final`)    |
| Abstract methods     | Yes                                   | Yes                                           |
| Concrete methods     | Yes                                   | Yes, using `default`/`static` methods         |
| State                | Can maintain object state             | Doesn't have normal instance state            |
| Main purpose         | Shared base/class implementation      | Contract/capability                           |
