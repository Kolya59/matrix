# matrix
Process-based matrix computing

Написать программу, которая реализовывала бы умножение матрицы на вектор при помощи нескольких процессов. 

Для  коммуникации процессов использовать механизм программных гнёзд (сокетов). 

Алгоритм:
Сначала главный процесс считывает исходные матрицу и вектор из файлов. Затем главный процесс создаёт несколько вспомогательных процессов и передаёт им несколько строк матрицы и исходный вектор столбец. После получения этих данных вспомогательные процессы должны перемножить соответствующие строки на вектор – столбец и передать полученные результаты главному процессу. Главный процесс должен вывести получившиеся результаты в файл.
 

Замечания:
•	Исходные матрица и вектор – столбец считываются из файла. 
Результат также записывается в файл. 
•	Необходимо сделать так, чтобы нагрузка на вспомогательные процессы была равномерной (то есть каждому из процессов досталось бы примерно одинаковое количество строк исходной матрицы)
