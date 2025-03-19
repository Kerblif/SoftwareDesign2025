using System;

namespace HSEFinance.Lib.Application.Import
{
    public abstract class FileImporterBase<T>
    {
        public T Import(string filePath)
        {
            if (string.IsNullOrEmpty(filePath))
                throw new ArgumentNullException(nameof(filePath), "File path cannot be null or empty.");

            var fileContent = LoadFileContent(filePath);
            return Parse(fileContent);
        }

        private string LoadFileContent(string filePath)
        {
            if (!File.Exists(filePath))
                throw new FileNotFoundException("File does not exist.", filePath);

            return File.ReadAllText(filePath);
        }

        protected abstract T Parse(string content);
    }
}