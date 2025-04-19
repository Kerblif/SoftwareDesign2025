using System.Text.Json;

namespace HSEFinance.Lib.Application.Export
{
    public class JsonFileExporter<T> : FileExporterBase<IEnumerable<T>>
    {
        protected override string Serialize(IEnumerable<T> data)
        {
            return JsonSerializer.Serialize(data, new JsonSerializerOptions
            {
                WriteIndented = true
            });
        }
    }
}