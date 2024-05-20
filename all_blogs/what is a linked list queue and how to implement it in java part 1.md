Title: What is a Linked List Queue and How to Implement it in Java, Part 1
Date: 20 May 2024
Author: Mohd. Sameer
Intro: In this post, I will explain what a Linked List Queue is and how to implement it in Java. Let's begin. In simple terms, a linked list queue is just a regular queue that can hold data in nodes or slots. The key characteristic is that the element added first will be removed first (FIFO - First In, First Out), similar to real-world queues, such as those in restaurants where the person who arrives first gets served first. In Big O notation, adding, removing, and peeking items are O(1) operations, meaning they are constant time operations regardless of the queue size. Now, let's implement it in Java.

## Hello World!

Assuming you have Java installed on your PC, create a new file named `Main.java`. First, let's create a "Hello World!" program.

Type the following code in it:

```java
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello World!");
    }
}
```
Now open your terminal in the same folder and type <kbd>javac Main.java</kbd> and press <kbd>enter</kbd> then type <kbd>java Main</kbd>. If you see `Hello World!` printed in your terminal, that means your program is running properly.

### Let's Keep going!

Now delete your **Hello World!** code and type the following code to create a node:

```java
public class Queue<T> {
    private static class QNode<T> {
        private T value;
        private QNode<T> next;

        public QNode(T value) {
            this.value = value;
            this.next = null;
        }
    }
}
```
The code above creates a class named `Queue` which is generic (`<T>`), meaning it can store any type of data (numbers, strings, etc.). Inside it, we have a nested class named `QNode` (short for Queue Node) that has two fields: `value` and `next`. The constructor of `QNode` initializes these fields.

Next, add the following:

```diff
public class Queue<T> {
    private static class QNode<T> {
        private T value;
        private QNode<T> next;

        public QNode(T value) {
            this.value = value;
            this.next = null;
        }
    }
+   
+   private QNode<T> head;
+   private QNode<T> tail;
+   private int length;
+   
+   public Queue() {
+       this.head = null;
+       this.tail = null;
+       this.length = 0;
+   }
}
```

In the code above, we define the `head`, `tail`, and `length` for the `Queue` class. `head` and `tail` will point to the first and last nodes in the queue, respectively, and `length` will store the number of nodes currently in the queue. We also have a constructor that initializes these fields.


### Let's Write More Code

```diff
public class Queue<T> {
    private static class QNode<T> {
        private T value;
        private QNode<T> next;

        public QNode(T value) {
            this.value = value;
            this.next = null;
        }
    }
    
    private QNode<T> head;
    private QNode<T> tail;
    private int length;
    
    public Queue() {
        this.head = null;
        this.tail = null;
        this.length = 0;
    }
+   
+   public int getLength() {
+       return length;
+   }
+   
+   public void enqueue(T item) {
+       QNode<T> node = new QNode<>(item);
+       length++;
+   
+       if (tail == null) {
+           head = tail = node;
+           return;
+       }
+   
+       tail.next = node;
+       tail = node;
+   }
}
```
We have added two more functions: `getLength()` and `enqueue(T item)`.

* `getLength()` is straightforward; it returns the length of the queue.
* `enqueue(T item)` takes a generic `T` parameter value, creates a new node with this value, and increments the queue length by 1. If the queue is empty (i.e., `tail` is `null`), it sets both head and tail to this new node. If the queue is not empty, it adds the new node to the end and updates the `tail`.

### Final Code Below

```java
public class Queue<T> {
    private static class QNode<T> {
        private T value;
        private QNode<T> next;

        public QNode(T value) {
            this.value = value;
            this.next = null;
        }
    }
    
    private QNode<T> head;
    private QNode<T> tail;
    private int length;
    
    public Queue() {
        this.head = null;
        this.tail = null;
        this.length = 0;
    }

    public int getLength() {
        return length;
    }

    public void enqueue(T item) {
        QNode<T> node = new QNode<>(item);
        length++;

        if (tail == null) {
            head = tail = node;
            return;
        }

        tail.next = node;
        tail = node;
    }
}
```
In the next part, we will implement the `dequeue` and `peek` functions to complete this Linked List Queue. Until next time! 👋
