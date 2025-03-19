using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;

namespace HSEFinance.Lib.Application.Import
{
    public class YamlFileImporter<T> : FileImporterBase<T>
    {
        protected override T Parse(string content)
        {
            var deserializer = new DeserializerBuilder()
                .WithNamingConvention(CamelCaseNamingConvention.Instance)
                .Build();

            return deserializer.Deserialize<T>(content);
        }
    }
}