using HSEFinance.Lib.Application.Export;
using HSEFinance.Lib.Application.Import;

namespace HSEFinance.Lib.Application.Facades
{
    public class ImportExportFacade<T>
    {
        private readonly Dictionary<string, FileExporterBase<IEnumerable<T>>> _exporters;
        private readonly Dictionary<string, FileImporterBase<IEnumerable<T>>> _importers;

        public ImportExportFacade()
        {
            // Регистрация экспортёров
            _exporters = new Dictionary<string, FileExporterBase<IEnumerable<T>>>
            {
                { "json", new JsonFileExporter<T>() },
                { "csv", new CsvFileExporter<T>() }
            };

            // Регистрация импортёров
            _importers = new Dictionary<string, FileImporterBase<IEnumerable<T>>>
            {
                { "json", new JsonFileImporter<T>() },
                { "csv", new CsvFileImporter<T>() }
            };
        }

        /// <summary>
        /// Экспортирует данные в указанный формат.
        /// </summary>
        /// <param name="data">Список объектов для экспорта</param>
        /// <param name="format">Формат файла (json, csv)</param>
        /// <param name="filePath">Путь к файлу</param>
        /// <exception cref="NotSupportedException">Если формат не поддерживается</exception>
        public void Export(IEnumerable<T> data, string format, string filePath)
        {
            if (!_exporters.TryGetValue(format.ToLower(), out var exporter))
            {
                var supportedFormats = string.Join(", ", _exporters.Keys);
                throw new NotSupportedException($"Формат \"{format}\" не поддерживается. Поддерживаемые форматы: {supportedFormats}");
            }

            try
            {
                exporter.Export(data, filePath);
            }
            catch (Exception ex)
            {
                throw new Exception($"Не удалось экспортировать данные в формат {format}. Ошибка: {ex.Message}", ex);
            }
        }

        /// <summary>
        /// Импортирует данные из файла указанного формата.
        /// </summary>
        /// <param name="format">Формат файла (json, csv)</param>
        /// <param name="filePath">Путь к файлу</param>
        /// <returns>Список объектов</returns>
        /// <exception cref="NotSupportedException">Если формат не поддерживается</exception>
        public IEnumerable<T>? Import(string format, string filePath)
        {
            if (!_importers.TryGetValue(format.ToLower(), out var importer))
            {
                var supportedFormats = string.Join(", ", _importers.Keys);
                throw new NotSupportedException($"Формат \"{format}\" не поддерживается. Поддерживаемые форматы: {supportedFormats}");
            }

            try
            {
                var data = importer.Import(filePath);
                return data;
            }
            catch (Exception ex)
            {
                throw new Exception($"Не удалось импортировать данные из файла {filePath} в формате {format}. Ошибка: {ex.Message}", ex);
            }
        }
    }
}