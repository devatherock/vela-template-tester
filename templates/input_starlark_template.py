def main(ctx):
  return {
    'version': '1',
    'steps': [
      {
        'name': 'build',
        'image': ctx["vars"]["image"],
        'commands': [
          'go build',
          'go test',
        ]
      },
    ],
}