# Project: Image Processing System + Work Stealing and Map Reduce

---
## **Project Description**

This project implements an **image processing system** using **Golang**, capable of applying effects like grayscale, sharpening, blurring, and edge detection on a set of images. The system supports both **serial** and two types of **parallel implementations**:

1. **Work Stealing Paradigm**
2. **Map Reduce**

Each parallelization method distributes tasks across threads, where each thread processes entire images.

---

## **Parallel Solution Description**

### **1. Work Stealing Paradigm**

- Utilizes a **Lock-Free Bounded Deque** as the data structure for local goroutine queues.
- A global goroutine queue is created from `effects.txt`, and tasks are distributed to threads using the `go` statement.
- Threads:
    - Push new tasks to the bottom of their own queue.
    - Steal tasks from the top of other threads' queues once they finish their own.
- Benefits:
    - Improved load balancing for mixed image sets.
    - Nearly linear speedup for smaller tasks due to better task distribution.


### **2. Map Reduce**

- Implements the classic MapReduce model:
    - **Mappers**: Read data from separate files and emit intermediate results via channels.
    - **Shuffler**: Groups intermediate results by the "Region" field.
    - **Reducers**: Process grouped data using goroutines, applying effects region-wise.
- Benefits:
    - Combines similar effects to improve load balancing and processing speed.
    - Parallelizes the reduce stage for faster execution.

---

## **How to Run**

1. Prepare the environment and dataset.
2. Use the following commands to test different implementations:
    - For Work Stealing on small images:

```bash
go run editor.go small ws 12
```

    - For Map Reduce on large images:

```bash
go run editor.go big mr 12
```

3. Run benchmark tests using:

```bash
sbatch benchmark.sh
```

4. Results will appear in the `editor` directory.

---

## **Challenges Faced**

- Balancing task granularity was critical for achieving speedup in both parallel implementations.
- Despite implementing work stealing, certain variables like idle threads and memory contention slowed down performance.
- In Map Reduce, mapping and shuffling stages were bottlenecks due to their sequential nature.

---

## **Results and Discussion**

### Speedup Plots

1. Speedup for small images processed using Map Reduce (sequential processing as reference).
2. Execution time comparison for Map Reduce with varying thread counts (1-thread execution as reference).

### Bottlenecks

- **Work Stealing**:
    - Generator remained sequential, limiting performance gains.
- **Map Reduce**:
    - Mapping and shuffling stages were not parallelized, creating significant bottlenecks.


### Hotspots

- **Work Stealing**:
    - Idle workers competing for tasks slowed down active threads.
- **Map Reduce**:
    - Uneven distribution of images across regions caused load imbalance.


### Speedup Limitations

- Overhead from communication and synchronization between threads limited linear scalability in both approaches.
- Optimal thread count depended on the number of images or regions being processed.

---

## **Comparison of Parallel Implementations**

| Feature | Work Stealing | Map Reduce |
| :-- | :-- | :-- |
| Task Distribution | Dynamic (work stealing) | Static (region-based grouping) |
| Load Balancing | Better for mixed image sizes | Limited by uneven region sizes |
| Speedup Consistency | Consistent across all image types | Dependent on number of regions |
| Bottlenecks | Idle workers \& generator | Mapping \& shuffling stages |

---

## **Conclusions**

- Work Stealing achieved better load balancing, especially for mixed image sets, resulting in consistent speedups.
- Map Reduce showed significant improvements for small images but was limited by load imbalance and sequential mapping/shuffling stages.

---

## **Future Improvements**

1. Parallelize bottleneck stages like mapping and shuffling in Map Reduce.
2. Optimize task granularity to improve performance consistency across all image sizes.
3. Explore hybrid approaches combining Work Stealing with Map Reduce for better scalability.

---

This project demonstrates how parallelization techniques like Work Stealing and Map Reduce can be applied effectively to computationally intensive tasks like image processing, highlighting both their strengths and limitations.

<div>⁂</div>

[^1]: https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/32275748/7191005f-f799-4a0e-95dd-7e4515a79e6e/Project-3-Report-Parallel.pdf

[^2]: https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/32275748/7191005f-f799-4a0e-95dd-7e4515a79e6e/Project-3-Report-Parallel.pdf
