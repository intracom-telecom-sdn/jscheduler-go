package test;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

public class Primes {
    public static void main(String[] args) {

        ExecutorService[] executors = {
                new ThreadPoolExecutor(
                        5, // core pool size
                        5, // max pool size
                        0, TimeUnit.MILLISECONDS, // keep alive
                        new LinkedBlockingQueue<>(1), // queue
                        new ThreadPoolExecutor.CallerRunsPolicy()),
                new ThreadPoolExecutor(
                        5, // core pool size
                        5, // max pool size
                        0, TimeUnit.MILLISECONDS, // keep alive
                        new LinkedBlockingQueue<>(1), // queue
                        new ThreadPoolExecutor.CallerRunsPolicy())
        }; // handler

        System.out.println("Main thread name " + Thread.currentThread().getName());
        for (int i = 0; i<2; i++) {
            for (int thread = 0; thread < 5; thread++) {
                executors[i].submit(() -> {
                    System.out.println("Entering thread " + Thread.currentThread().getName() + " #" + Thread.currentThread().getId());

                    long timeStart = System.currentTimeMillis();
                    long timeEnd = 0L;
                    double timeDiff = 0.0;
                    double throughput = 0.0;
                    for (long job = 1; job <= 1000000000; job++) {
                        long count = 0;
                        for (long k=3; k<=10000L; k++) {
                            boolean isPrime = true;
                            for (long j=2; j<=k/2 && isPrime; j++) {
                                isPrime = k % j > 0;
                            }
                            if (isPrime) {
                                count++;
                            }
                        }
                        if(job % 100 == 0) {
                            timeEnd = System.currentTimeMillis();
                            timeDiff = (timeEnd - timeStart) / 1000.0;
                            throughput = 100.0 / timeDiff;
                            System.out.printf("  %s throughput: %.2f jobs/sec\n", Thread.currentThread().getName(), throughput);
                            timeStart = System.currentTimeMillis();
                        } 
                    }
                });
            }
            executors[i].shutdown();
        }

    }

}
