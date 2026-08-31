#include <vector>
 #include <set>
 #include <map>
 #include <algorithm>
 #include <memory>
 #include <functional>

 using namespace std;

 class Product {
 public:
     int id;
     int count;

     Product(int id, int count) : id(id), count(count) {}

     bool operator<(const Product& other) const {
         id < other.id;
     }
 };

 // Контекст для передачи данных через пайплайн
 struct PipelineContext {
     map<int, set<int>> userItems;
     map<int, int> itemUserCount;
     vector<Product> result;
 };

 // Базовый класс шага пайплайна
 class PipelineStep {
 public:
     virtual ~PipelineStep() = default;
     virtual PipelineContext execute(const vector<vector<int>>& input, PipelineContext context) = 0;
 };

 // Шаг 1: Сбор данных о пользователях
 class UserDataCollectorStep : public PipelineStep {
 public:
     PipelineContext execute(const vector<vector<int>>& input, PipelineContext context) override {
         //ваш код
     }
 };

 // Шаг 2: Подсчет пользователей на товар
 class ItemUserCounterStep : public PipelineStep {
 public:
     PipelineContext execute(const vector<vector<int>>& input, PipelineContext context) override {
         //ваш код
     }
 };

 // Шаг 3: Фильтрация и сортировка товаров
 class ItemFilterStep : public PipelineStep {
 public:
     PipelineContext execute(const vector<vector<int>>& input, PipelineContext context) override {
        //ваш код
     }
 };

 // Пайплайн обработки
 class ProcessingPipeline {
 private:
     vector<shared_ptr<PipelineStep>> steps;

 public:
     void addStep(shared_ptr<PipelineStep> step) {
         steps.push_back(step);
     }

     PipelineContext execute(const vector<vector<int>>& input) {
        //ваш код
     }
 };

 class ProcessItems {
 public:
 //функция main дополнительно создает тестовый конвейер для одного товара
 //если комплементарных товаров нет возвращается Product(0,0)
     vector<Product> getComplementaryProducts(const vector<vector<int>>& input) {
         //ваш код
     }
 };