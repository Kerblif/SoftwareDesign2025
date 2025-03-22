using System.Text.Json;

namespace HSEFinance.Lib.Application.Import
{
    public class JsonFileImporter<T> : FileImporterBase<IEnumerable<T>>
    {
        protected override IEnumerable<T> Parse(string content)
        {
            try
            {
                return JsonSerializer.Deserialize<IEnumerable<T>>(content)
                       ?? throw new InvalidOperationException("Failed to deserialize JSON content.");
            }
            catch (JsonException ex)
            {
                throw new InvalidOperationException("Invalid JSON format.", ex);
            }
        }
    }
}