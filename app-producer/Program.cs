using Bogus;
// See https://aka.ms/new-console-template for more information
Console.WriteLine("Hello, World!");


var faker = new Faker<Entry>()
    .CustomInstantiator(f => new Entry(Guid.NewGuid(), DateTime.Now, f.Lorem.Sentence()));

faker.Generate(10).ForEach(e => Console.WriteLine(e));