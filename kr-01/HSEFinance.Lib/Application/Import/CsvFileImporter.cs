using CsvHelper;
using System.Globalization;

namespace HSEFinance.Lib.Application.Import
{
    public class CsvFileImporter<T> : FileImporterBase<IEnumerable<T>>
    {
        protected override IEnumerable<T> Parse(string content)
        {
            using var reader = new StringReader(content);
            using var csv = new CsvReader(reader, CultureInfo.InvariantCulture);

            return csv.GetRecords<T>().ToList();
        }
    }
}