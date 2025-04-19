using CsvHelper;
using CsvHelper.Configuration;
using System.Globalization;

namespace HSEFinance.Lib.Application.Export
{
    public class CsvFileExporter<T> : FileExporterBase<IEnumerable<T>>
    {
        protected override string Serialize(IEnumerable<T> data)
        {
            using var writer = new StringWriter();
            var csvConfig = new CsvConfiguration(CultureInfo.InvariantCulture)
            {
                Delimiter = ",",
                Encoding = System.Text.Encoding.UTF8,
            };
            using var csv = new CsvWriter(writer, csvConfig);
            
            csv.WriteRecords(data);
            return writer.ToString();
        }
    }
}