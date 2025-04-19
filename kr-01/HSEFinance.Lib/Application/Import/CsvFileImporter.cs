using CsvHelper;
using System.Globalization;
using HSEFinance.Lib.Domain.Entities;

namespace HSEFinance.Lib.Application.Import
{
    public class CsvFileImporter<T> : FileImporterBase<IEnumerable<T>>
    {
        protected override IEnumerable<T> Parse(string content)
        {
            using var reader = new StringReader(content);
            var csvConfig = new CsvHelper.Configuration.CsvConfiguration(CultureInfo.InvariantCulture)
            {
                HeaderValidated = null,
                MissingFieldFound = null,
                PrepareHeaderForMatch = args => args.Header.ToLower()
            };

            using var csv = new CsvReader(reader, csvConfig);
            
            return csv.GetRecords<T>().ToList();
        }
    }
}