using System.Text.Json;

namespace HSEFinance.Lib.Application.Export
{
    public class JsonFileExporter<T> : FileExporterBase<T>
    {
        protected override string Serialize(T data)
        {
            return JsonSerializer.Serialize(data, new JsonSerializerOptions
            {
                WriteIndented = true
            });
        }
    }
}