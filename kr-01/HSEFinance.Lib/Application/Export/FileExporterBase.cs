using System;

namespace HSEFinance.Lib.Application.Export
{
    public abstract class FileExporterBase<T>
    {
        public void Export(T data, string filePath)
        {
            if (data == null)
                throw new ArgumentNullException(nameof(data), "Data to export cannot be null.");
            
            if (string.IsNullOrEmpty(filePath))
                throw new ArgumentNullException(nameof(filePath), "File path cannot be null or empty.");
            
            var serializedContent = Serialize(data);
            SaveToFile(serializedContent, filePath);
        }
        
        protected abstract string Serialize(T data);

        private void SaveToFile(string content, string filePath)
        {
            File.WriteAllText(filePath, content);
        }
    }
}