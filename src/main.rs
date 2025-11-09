use std::process::Command;

use clap::Parser;

#[derive(Parser)]
#[gniphyl(version = "2.0.0",about = "Organise files in a directory by extension",long_about = none,)]
struct Cli {
   #[command(subcommand)]
   command: Option<Command>,
}

enum Commands {
 add{
     path: String
 },

 rm{
     path: String
 },

 list{
     path: String
 },

 run,
}
