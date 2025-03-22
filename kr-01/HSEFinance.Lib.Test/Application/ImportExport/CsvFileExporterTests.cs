using System.Globalization;
using CsvHelper;
using CsvHelper.Configuration;
using HSEFinance.Lib.Application.Export;
using HSEFinance.Lib.Application.Facades;
using HSEFinance.Lib.Application.Import;
using Xunit;

namespace HSEFinance.Lib.Test.Application.Export
{
    public class CsvFileExporterTests
    {
        ImportExportFacade<TestData> fasade = new ImportExportFacade<TestData>();
        
        [Fact]
        public void ImportExport_WithValidData_ReturnsCsvString()
        {
            var data = new List<TestData>
            {
                new TestData { Id = 1, Name = "Test 1" },
                new TestData { Id = 2, Name = "Test 2" }
            };

            var tempFilePath = Path.GetTempFileName();
            fasade.Export(data, "csv", tempFilePath);
            
            Assert.True(File.Exists(tempFilePath));
            
            var parsedData = (fasade.Import("csv", tempFilePath) ?? Array.Empty<TestData>()).ToList();
            Assert.NotNull(parsedData);
            Assert.Equal(2, parsedData.Count);
            Assert.Equal(1, parsedData[0].Id);
            Assert.Equal("Test 1", parsedData[0].Name);
            Assert.Equal(2, parsedData[1].Id);
            Assert.Equal("Test 2", parsedData[1].Name);
        }

        [Fact]
        public void ImportExport_WithEmptyData_ReturnsOnlyHeader()
        {
            var data = new List<TestData>();

            var tempFilePath = Path.GetTempFileName();
            fasade.Export(data, "csv", tempFilePath);
            
            Assert.True(File.Exists(tempFilePath));
            
            var parsedData = fasade.Import("csv", tempFilePath);
            Assert.NotNull(parsedData);
            Assert.Empty(parsedData);
        }

        [Fact]
        public void Export_WithInvalidData_ThrowsException()
        {
            Assert.Throws<Exception>(() => fasade.Export(null, "json", "./tmp"));
        }

        private class TestData
        {
            public int Id { get; set; }
            public string Name { get; set; } = string.Empty;
        }
    }
}