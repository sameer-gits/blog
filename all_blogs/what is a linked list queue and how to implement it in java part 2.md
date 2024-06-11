Title: What is a Linked List Queue and How to Implement it in Java, Part 2
Date: 11 Jun 2024
Author: Mohd. Sameer
Intro: Welcome back! In our previous post, we explored the concept of a Linked List Queue and began its implementation in Java. We established fundamental functionalities like `getLength` and `enqueue`. Now, let’s expand our queue’s capabilities by implementing the `peek` and `dequeue` functions. These methods are crucial for a fully functional queue, allowing us to view and remove elements in a First-In-First-Out (FIFO) manner. To recap, a linked list queue is a data structure where elements are stored in nodes, with each node pointing to the next, thus creating a sequence. In this structure, the first element added is the first to be removed, mimicking a real-world queue.

## Previous code

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

## peek

Open your previous file and Type the following code in it:

```java
public T peek() {
    if (head == null) {
        return null;
    }
    return head.value;
}
```

Above code will return the value if it's there and if there is no value it will just return null, Let's write more code.

## dequeue

```java
public T dequeue() {
    if (head == null) {
        return null;
    }

    length--;
    T item = head.value;
    head = head.next;

    if (head == null) {
        tail = null;
    }
    return item;
}
```

Above code is for removing an item from the head of the Queue if head is null it return `null` and else it decrement the length of the Queue and then assigns `head.value` to a temporary `item` for returning it in future and make `head = head.next` if then head is `null` then it also make tail `null` and return `item`.

## Complete Code

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

    public T peek() {
        if (head == null) {
            return null;
        }
        return head.value;
    }

    public T dequeue() {
        if (head == null) {
            return null;
        }

        length--;
        T item = head.value;
        head = head.next;

        if (head == null) {
            tail = null;
        }
        return item;
    }
}
```

And that's it! We've successfully implemented a Linked List Queue in Java, complete with methods to enqueue, dequeue, peek, and get the length of the queue. This implementation ensures efficient and straightforward operations for managing a queue.

In summary, we've covered the essential methods needed for a functional Linked List Queue. With enqueue, dequeue, and peek, we can add, remove, and inspect elements in our queue efficiently. Until next time! 👋