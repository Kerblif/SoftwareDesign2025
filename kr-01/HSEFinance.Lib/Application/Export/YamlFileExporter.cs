using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;

namespace HSEFinance.Lib.Application.Export
{
    public class YamlFileExporter<T> : FileExporterBase<T>
    {
        protected override string Serialize(T data)
        {
            var serializer = new SerializerBuilder()
                .WithNamingConvention(CamelCaseNamingConvention.Instance)
                .Build();

            return serializer.Serialize(data);
        }
    }
}