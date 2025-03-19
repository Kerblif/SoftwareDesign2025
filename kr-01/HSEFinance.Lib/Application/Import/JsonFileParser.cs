using System.Text.Json;

namespace HSEFinance.Lib.Application.Import
{
    public class JsonFileImporter<T> : FileImporterBase<T>
    {
        protected override T Parse(string content)
        {
            try
            {
                return JsonSerializer.Deserialize<T>(content)
                       ?? throw new InvalidOperationException("Failed to deserialize JSON content.");
            }
            catch (JsonException ex)
            {
                throw new InvalidOperationException("Invalid JSON format.", ex);
            }
        }
    }
}